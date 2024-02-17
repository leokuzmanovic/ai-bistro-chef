package services

import (
	"context"
	"fmt"
	"os"

	"github.com/leokuzmanovic/ai-bistro-chef/internal/configuration"
	errs "github.com/leokuzmanovic/ai-bistro-chef/internal/errors"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const (
	assistant_Function_Describe_Image             = "describe_image"
	assistant_Function_Describe_Image_Description = `"Describe an image user uploaded, understand the individual items (ingredients) in the image, 
	and use it as an input for constructing an appropriate prompt that will be later used to generate a cooking recipe."`
)

type openaiAssistantBuilder interface {
	initAssistant() (*openai.Assistant, error)
}

type openaiAssistantBuilderImpl struct {
	client           *openai.Client
	assistantConfig  *configuration.OpenAiAssistantConfig
	localRecipesPath string
}

func newOpenaiAssistantBuilderImpl(client *openai.Client, localRecipesPath string, assistantConfig *configuration.OpenAiAssistantConfig) *openaiAssistantBuilderImpl {
	return &openaiAssistantBuilderImpl{
		client:           client,
		assistantConfig:  assistantConfig,
		localRecipesPath: localRecipesPath,
	}
}

func (s *openaiAssistantBuilderImpl) initAssistant() (*openai.Assistant, error) {
	assistant, err := s.loadAssistant()
	if err != nil || assistant == nil {
		return s.createNew()
	} else {
		fmt.Println("Assistant already exists")
		err = s.updateAssistant(assistant)
		if err != nil {
			return assistant, err
		}
		return assistant, nil
	}
}

func (s *openaiAssistantBuilderImpl) createNew() (*openai.Assistant, error) {
	assistant, err := s.createNewAssistant()
	if err != nil {
		return assistant, err
	}
	return assistant, nil
}

func (s *openaiAssistantBuilderImpl) loadAssistant() (*openai.Assistant, error) {
	limit := 20
	order := "desc"
	var before *string = nil
	var after *string = nil

	for {
		assistantsListResponse, err := s.client.ListAssistants(context.Background(), &limit, &order, before, after)
		if err != nil && errs.IsOpenAINotFoundError(err) {
			return nil, err
		} else if err != nil {
			fmt.Println("Error while listing assistants!")
			return nil, err
		}

		if assistantsListResponse.Assistants == nil || len(assistantsListResponse.Assistants) == 0 {
			return nil, nil
		}

		before = assistantsListResponse.LastID
		for _, assistant := range assistantsListResponse.Assistants {
			if *assistant.Name == *s.assistantConfig.GetAssistantName() {
				return &assistant, nil
			}
		}

		if !assistantsListResponse.HasMore {
			return nil, nil
		}
	}
}

func (s *openaiAssistantBuilderImpl) updateAssistant(assistant *openai.Assistant) error {
	assistantFileIds, updated := s.checkAssistentFiles(assistant.ID, assistant.FileIDs)

	if updated {
		fmt.Println("Updating assistant...")
		_, err := s.client.ModifyAssistant(context.Background(), assistant.ID,
			openai.AssistantRequest{
				FileIDs: assistantFileIds,
			})

		if err != nil {
			fmt.Println("Error while updating assistant!")
			return err
		}
		fmt.Println("Assistant updated!")
	}
	return nil
}

func (s *openaiAssistantBuilderImpl) createNewAssistant() (*openai.Assistant, error) {
	fmt.Println("Creating assistant")
	assistant, err := s.prepareNewAssistant()
	if err != nil {
		return nil, err
	}
	fmt.Println("Assistant created")
	return assistant, nil
}

func (s *openaiAssistantBuilderImpl) checkAssistentFiles(assistantId string, assistantFileIds []string) ([]string, bool) {
	openAIFiles := make([]openai.File, 0)
	for _, fileId := range assistantFileIds {
		openAIFile, err := s.client.GetFile(context.Background(), fileId)
		if err != nil {
			fmt.Println("Error while getting file!")
			panic(err)
		}
		openAIFiles = append(openAIFiles, openAIFile)
	}
	return s.prepareAssistentFiles(openAIFiles)
}

func (s *openaiAssistantBuilderImpl) prepareNewAssistant() (*openai.Assistant, error) {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"user_message": {
				Type:        jsonschema.String,
				Description: "Message the user sent along with the image, without the URL itself",
			},
			"image": {
				Type:        jsonschema.String,
				Description: "URL of the image user provided",
			},
		},
		Required: []string{"user_message", "image"},
	}

	tools := []openai.AssistantTool{
		{
			Type: openai.AssistantToolTypeCodeInterpreter,
		},
		{
			Type: openai.AssistantToolTypeRetrieval,
		},
		{
			Type: openai.AssistantToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        assistant_Function_Describe_Image,
				Description: assistant_Function_Describe_Image_Description,
				Parameters:  params,
			},
		},
	}

	// make sure all assistant files are synced with openai
	filesFromOpenAI, err := s.client.ListFiles(context.Background())
	if err != nil {
		fmt.Println("Error while listing files!")
		return nil, err
	}

	assistentFiles, _ := s.prepareAssistentFiles(filesFromOpenAI.Files)

	assistant, err := s.client.CreateAssistant(context.Background(),
		openai.AssistantRequest{
			Name:         s.assistantConfig.GetAssistantName(),
			Model:        *s.assistantConfig.GetAssistantGptModel(),
			Description:  s.assistantConfig.GetAssistantDescription(),
			Instructions: s.assistantConfig.GetAssistantInstructions(),
			Tools:        tools,
			FileIDs:      assistentFiles,
		})
	return &assistant, err
}

func (s *openaiAssistantBuilderImpl) prepareAssistentFiles(filesFromOpenAI []openai.File) ([]string, bool) {
	recipeFileIds, updated := s.synchroniseRecipes(filesFromOpenAI)
	return recipeFileIds, updated
}

func (s *openaiAssistantBuilderImpl) synchroniseRecipes(filesFromOpenAI []openai.File) ([]string, bool) {
	// NOTE: ideally we would use the file hash to check if it needs to be replaced, but openai does allow assisstant files to be downloaded
	filesToUploadMap := make(map[string]int)       // [filename][byte size]
	fileIdsToDeleteFromOpenAI := make([]string, 0) // fileIds
	assistantFileIds := make([]string, 0)          // fileIds
	updateNeeded := false
	s.readAllLocalRecipes(s.localRecipesPath, filesToUploadMap)

	// check if remote file exists in local files map and if hashes are different to decide if file should be replaced
	for _, openAIFile := range filesFromOpenAI {
		openAIFileSize := s.getOpenAIFileByteSize(openAIFile)
		if localFileSize, ok := filesToUploadMap[openAIFile.FileName]; ok {
			if localFileSize != openAIFileSize {
				// remote file needs to be replaced, first add its current id to be deleted (local one will be uploaded later)
				fileIdsToDeleteFromOpenAI = append(fileIdsToDeleteFromOpenAI, openAIFile.ID)
			} else {
				// local file is already uploaded to openai, remove it from map
				delete(filesToUploadMap, openAIFile.FileName)
				assistantFileIds = append(assistantFileIds, openAIFile.ID)
			}
		} else {
			// local file list does not include this remote file anymore, it needs to be deleted
			fileIdsToDeleteFromOpenAI = append(fileIdsToDeleteFromOpenAI, openAIFile.ID)
		}
	}

	// delete files from openai if needed
	if len(fileIdsToDeleteFromOpenAI) > 0 {
		updateNeeded = true
		s.deleteOpenAIFiles(fileIdsToDeleteFromOpenAI)
	}

	// upload local files to openai
	if len(filesToUploadMap) > 0 {
		updateNeeded = true
		assistantFileIds = append(assistantFileIds, s.uploadFilesToOpenAI(s.localRecipesPath, filesToUploadMap)...)
	}
	return assistantFileIds, updateNeeded
}

func (s *openaiAssistantBuilderImpl) uploadFilesToOpenAI(localRecipesPath string, filesToUploadMap map[string]int) []string {
	fmt.Println("Uploading files to openai...")
	uploadedFileIds := make([]string, 0)

	for filename := range filesToUploadMap {
		data, err := os.ReadFile(s.localRecipesPath + "/" + filename)
		if err != nil {
			fmt.Println("Error while reading files!")
			panic(err)
		}
		file, err := s.client.CreateFileBytes(context.Background(), openai.FileBytesRequest{
			Name:    filename,
			Bytes:   data,
			Purpose: openai.PurposeAssistants,
		})
		if err != nil {
			fmt.Println("Error while uploading files!")
			panic(err)
		}
		uploadedFileIds = append(uploadedFileIds, file.ID)
	}
	return uploadedFileIds
}

func (s *openaiAssistantBuilderImpl) deleteOpenAIFiles(fileIdsToDeleteFromOpenAI []string) {
	fmt.Println("Deleting files from openai...")
	for _, fileId := range fileIdsToDeleteFromOpenAI {
		err := s.client.DeleteFile(context.Background(), fileId)
		if err != nil {
			fmt.Println("Error while deleting files!")
			panic(err)
		}
	}
	fmt.Println("Files deleted!")
}

func (s *openaiAssistantBuilderImpl) getOpenAIFileByteSize(openAIFile openai.File) int {
	file, err := s.client.GetFile(context.Background(), openAIFile.ID)
	if err != nil {
		fmt.Println("Error while getting file content!")
		panic(err)
	}
	return file.Bytes
}

func (*openaiAssistantBuilderImpl) readAllLocalRecipes(localRecipesPath string, filesToUploadMap map[string]int) {
	files, err := os.ReadDir(localRecipesPath)
	if err != nil {
		fmt.Println("Error while reading recipes folder!")
		panic(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(localRecipesPath + "/" + file.Name())
		if err != nil {
			panic(err)
		}

		filesToUploadMap[file.Name()] = len(data)
	}
}
