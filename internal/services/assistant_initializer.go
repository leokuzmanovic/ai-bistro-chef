package services

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type AssistantInitializer interface {
}

func (s *AssistantServiceImpl) PrepareAssistant() error {
	assistant, err := s.fetchAssistant()
	if err != nil || assistant == nil {
		return s.createNewAssistant()
	} else {
		fmt.Println("Assistant already exists!")
		return s.updateAssistant(assistant)
	}
}

func (s *AssistantServiceImpl) fetchAssistant() (*openai.Assistant, error) {
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

func (s *AssistantServiceImpl) updateAssistant(assistant *openai.Assistant) error {
	assistantFileIds := s.checkAssistentFiles(assistant.ID, assistant.FileIDs)

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
	return nil
}

func (s *AssistantServiceImpl) createNewAssistant() error {
	fmt.Println("Creating assistant...")
	_, err := s.prepareNewAssistant()
	if err != nil {
		return err
	}
	fmt.Println("Assistant created!")
	return nil
}

func (s *AssistantServiceImpl) checkAssistentFiles(assistantId string, assistantFileIds []string) []string {
	openAIFiles := make([]openai.File, 0)
	for _, fileId := range assistantFileIds {
		openAIFile, err := s.client.GetFile(context.Background(), fileId)
		if err != nil {
			fmt.Println("Error while getting file!")
			panic(err)
		}
		openAIFiles = append(openAIFiles, openAIFile)
	}
	recipeFileIds := s.prepareAssistentFiles(openAIFiles)
	return recipeFileIds
}

func (s *AssistantServiceImpl) prepareNewAssistant() (openai.Assistant, error) {
	tools := []openai.AssistantTool{
		{
			Type: openai.AssistantToolTypeCodeInterpreter,
		},
		{
			Type: openai.AssistantToolTypeRetrieval,
		},
	}

	// make sure all assistant files are synced with openai
	filesFromOpenAI, err := s.client.ListFiles(context.Background())
	if err != nil {
		fmt.Println("Error while listing files!")
		return openai.Assistant{}, err
	}

	assistant, err := s.client.CreateAssistant(context.Background(),
		openai.AssistantRequest{
			Name:         s.assistantConfig.GetAssistantName(),
			Model:        CHATGPT_MODEL,
			Description:  s.assistantConfig.GetAssistantDescription(),
			Instructions: s.assistantConfig.GetAssistantInstructions(),
			Tools:        tools,
			FileIDs:      s.prepareAssistentFiles(filesFromOpenAI.Files),
		})
	return assistant, err
}

func (s *AssistantServiceImpl) prepareAssistentFiles(filesFromOpenAI []openai.File) []string {
	recipeFileIds := s.synchroniseRecipes(filesFromOpenAI)
	return recipeFileIds
}

func (s *AssistantServiceImpl) synchroniseRecipes(filesFromOpenAI []openai.File) []string {
	// NOTE: ideally we would use hash of a file to check if it needs to be replaced, but openai does allow assisstant files to be downloaded
	filesToUploadMap := make(map[string]int)       // [filename][byte size]
	fileIdsToDeleteFromOpenAI := make([]string, 0) // fileIds
	assistantFileIds := make([]string, 0)          // fileIds

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
		s.deleteOpenAIFiles(fileIdsToDeleteFromOpenAI)
	}

	// upload local files to openai
	if len(filesToUploadMap) > 0 {
		assistantFileIds = append(assistantFileIds, s.uploadFilesToOpenAI(s.localRecipesPath, filesToUploadMap)...)
	}
	return assistantFileIds
}

func (s *AssistantServiceImpl) uploadFilesToOpenAI(localRecipesPath string, filesToUploadMap map[string]int) []string {
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

func (s *AssistantServiceImpl) deleteOpenAIFiles(fileIdsToDeleteFromOpenAI []string) {
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

func (s *AssistantServiceImpl) getOpenAIFileByteSize(openAIFile openai.File) int {
	file, err := s.client.GetFile(context.Background(), openAIFile.ID)
	if err != nil {
		fmt.Println("Error while getting file content!")
		panic(err)
	}
	return file.Bytes
}

func (*AssistantServiceImpl) readAllLocalRecipes(localRecipesPath string, filesToUploadMap map[string]int) {
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
