# AI Bistro Chef

AI Bistro Chef is an OpenAI Assistent and Vision API demo project designed to create a virtual assistant chef that can suggest meals based on a set of ingredients provided by the user, either through text or images (of ingredients). This assistant, or "chef," is trained on predefined recipes representing its knowledge base of favorite foods it can prepare. Beyond leveraging its internal knowledge, the chef can also browse the internet to find new meal suggestions, ensuring users always have something exciting to try.

## Functional Requirements

- **Assistant Creation:** Automatically creates an assistant (chef) if not already existing.
- **Training:** Trains the assistant with predefined recipes to form a knowledge base of meals the chef can suggest.
- **Meal Suggestions:** Prompts the chef with text or an image of ingredients to suggest a meal. The chef responds with a recipe from its knowledge base.
- **Internet Browsing:** If no meals can be recommended from the knowledge base, the chef can browse the internet to find suitable recipes.

### API Endpoints

This project provides a set of RESTful API endpoints to interact with the virtual assistant chef. Below are the endpoints available:

- **GET `/conversations`**
  
  Lists all conversations. This endpoint allows the retrieval of all ongoing conversations between users and the AI chef.

- **POST `/conversations/:conversationId/message`**
  
  Sends a message within a specific conversation. This endpoint is used to send either a text or an image of ingredients to the chef in an existing conversation, identified by `:conversationId`.

- **GET `/conversations/:conversationId/message/:messageId`**
  
  Retrieves a specific message's response by the chef. This endpoint allows fetching the chef's response to a specific message, identified by `:messageId`, within a conversation identified by `:conversationId`.


## Development Todos

- Enforce authentication to prevent tampering with threads that do not belong to users calling endpoints.
- Implement periodic clearing of local cache holding thread run results or adopt a centralized caching system like Redis with TTL to manage cache efficiently.
- Develop stream completions for assistant thread runs to enable streaming of results as tokens become available.

## Future Enhancements

- **Testing:** Integrate comprehensive test suites to ensure reliability and stability.
- **Configuration:** Externalize configuration and secrets management to enhance security and flexibility.


### Prompt engineering - example
persona - context - task - format - exemplar - tone

persona: your are a senior product marketing manager at Apple

context: and you have just unveiled the latest Apple product in collaboration with Tesla, the Apple Car, and received 12,000 pre-orders, which is 200% higher than target.

task: Write an...

format: email 

exemplar: This email should contain tl;dr (too long, didn't read) section, project background (why this product came into existance), business results section (quantifiable business metrics), and end with a section thanking the product and engineering teams.

tone: Use clear and concise language and write in a confident yet friendly tone.

