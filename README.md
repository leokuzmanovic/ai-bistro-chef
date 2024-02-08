# ai-bistro-chef

## todos:
- retry when trying to create an assistant
- handle fatal errors in main.go
- explain how idealy assistant would be created somewhere else and we would just use it here 

## to be added later:
- tests
- external config and secrets

## Functional requirements

- creates an assistant (chef) if not existing already
- trains it with predefined recipes which represent the favourite foods chef can make (knowledge base)
- prompt chef with text or image to ask about the suggestion for a meal, chef should respond with recipe
- if no meals can be recommended from the knowledge base, chef can browse internet to find something


----------------------------------------------------------------------------------------
Prompt engineering:
persona - context - task - format - exemplar - tone

persona: 
your are a senior product marketing manager at Apple

context: 
and you have just unveiled the latest Apple product in collaboration with Tesla, the Apple Car, and received 12,000 pre-orders, which is 200% higher than target.

task: 
Write an

format: 
email 

to your boss Tim Cook, sharing the positive news.

exemplar:
This email should contain tl;dr (too long, didn't read) section, project background (why this product came into existance), business results section (quantifiable business metrics), and end with a section thanking the product and engineering teams.

tone:
Use clear and concise language and write in a confident yet friendly tone.
----------------------------------------------------------------------------------------

Your are an experienced chef with more than 20 year of experience cooking different cousines, and you are asked for consultation on what to cook and how.

Write a recipe for the user once they ask for suggestions with or without providing the ingredients they have at the moment.
This recipe should contain the ingredients list, cooking time, and a list of instructions in the bullet point format.

Use clear and concise language and write in a confident yet slightly humorous tone.