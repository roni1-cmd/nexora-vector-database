import argparse
import os
from typing import List

from dotenv import load_dotenv
from openai import OpenAI
from openai.types.chat import ChatCompletionMessageParam
import chromadb


def build_prompt(query: str, context: List[str]) -> List[ChatCompletionMessageParam]:
    """
    Builds a prompt for the LLM. #

    This function builds a prompt for the LLM. It takes the original query,
    and the returned context, and asks the model to answer the question based only
    on what's in the context, not what's in its weights.

    More information: https://platform.openai.com/docs/guides/chat/introduction

    Args:
    query (str): The original query.
    context (List[str]): The context of the query, returned by embedding search.

    Returns:
    A prompt for the LLM (List[ChatCompletionMessageParam]).
    """

    system: ChatCompletionMessageParam = {
        "role": "system",
        "content": "I am going to ask you a question, which I would like you to answer"
        "based only on the provided context, and not any other information."
        "If there is not enough information in the context to answer the question,"
        'say "I am not sure", then try to make a guess.'
        "Break your answer up into nicely readable paragraphs.",
    }
    user: ChatCompletionMessageParam = {
        "role": "user",
        "content": f"The question is {query}. Here is all the context you have:"
        f'{(" ").join(context)}',
    }

    return [system, user]


def get_chatGPT_response(query: str, context: List[str], model_name: str) -> str:
    """
    Queries the GPT API to get a response to the question.

    Args:
    query (str): The original query.
    context (List[str]): The context of the query, returned by embedding search.

    Returns:
    A response to the question.
    """
    client = OpenAI()
    response = client.chat.completions.create(
        model=model_name,
        messages=build_prompt(query, context),
    )

    return response.choices[0].message.content  # type: ignore


def main(
    collection_name: str = "documents_collection", persist_directory: str = "."
) -> None:
    load_dotenv()

    # Check if the OPENAI_API_KEY environment variable is set.
    if not os.getenv("OPENAI_API_KEY"):
        print(
            "Please enter your OpenAI API Key. You can get it from https://platform.openai.com/account/api-keys\n"
        )
        return

    # Ask what model to use
    model_name = "gpt-4o-mini"
    change_model = input(f"This program is using the GPT-4o-mini model. Do you want to use a different model? (y/n): ")
    if change_model == "y":
        model_name = input("Please enter your desired model: ")

    # Instantiate a persistent chroma client in the persist_directory.
    # This will automatically load any previously saved collections.
    # Learn more at docs.trychroma.com
    client = chromadb.PersistentClient(path=persist_directory)

    # Get the collection.
    collection = client.get_collection(name=collection_name)

    # We use a simple input loop.
    while True:
        # Get the user's query
        query = input("Query: ")
        if len(query) == 0:
            print("Please enter a question. Ctrl+C to Quit.\n")
            continue
        print(f"\nThinking using {model_name}...\n")

        # Query the collection to get the 5 most relevant results
        results = collection.query(
            query_texts=[query], n_results=5, include=["documents", "metadatas"]
        )

        sources = "\n".join(
            [
                f"{result['filename']}: line {result['line_number']}"
                for result in results["metadatas"][0]  # type: ignore
            ]
        )

        # Get the response from GPT
        response = get_chatGPT_response(query, results["documents"][0], model_name)  # type: ignore

        # Output, with sources
        print(response)
        print("\n")
        print(f"Source documents:\n{sources}")
        print("\n")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Load documents from a directory into a Chroma collection"
    )

    parser.add_argument(
        "--persist_directory",
        type=str,
        default="chroma_storage",
        help="The directory where you want to store the Chroma collection",
    )
    parser.add_argument(
        "--collection_name",
        type=str,
        default="documents_collection",
        help="The name of the Chroma collection",
    )

    # Parse arguments
    args = parser.parse_args()

    main(
        collection_name=args.collection_name,
        persist_directory=args.persist_directory,
    )
