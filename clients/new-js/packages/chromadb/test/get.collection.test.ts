import { beforeEach, describe, expect, test } from "@jest/globals";
import { ChromaClient, CloudClient } from "../src";
import { DOCUMENTS, EMBEDDINGS, IDS, METADATAS } from "./utils/data";
import { DefaultEmbeddingFunction } from "@chroma-core/default-embed";

describe("get collections", () => {
  // connects to the unauthenticated chroma instance started in
  // the global jest setup file.
  const client = new ChromaClient({
    path: process.env.DEFAULT_CHROMA_INSTANCE_URL,
  });

  beforeEach(async () => {
    await client.reset();
  });

  test("it should get documents from a collection", async () => {
    const collection = await client.createCollection({ name: "test" });
    await collection.add({
      ids: IDS,
      embeddings: EMBEDDINGS,
      metadatas: METADATAS,
    });
    const results = await collection.get({ ids: ["test1"] });
    expect(results?.ids).toHaveLength(1);
    expect(["test1"]).toEqual(expect.arrayContaining(results.ids));
    expect(["test2"]).not.toEqual(expect.arrayContaining(results.ids));
    expect(results.include).toEqual(
      expect.arrayContaining(["metadatas", "documents"]),
    );

    const results2 = await collection.get({
      where: { test: "test1" },
    });
    expect(results2?.ids).toHaveLength(1);
    expect(["test1"]).toEqual(expect.arrayContaining(results2.ids));
  });

  test("wrong code returns an error", async () => {
    const collection = await client.createCollection({ name: "test" });
    await collection.add({
      ids: IDS,
      embeddings: EMBEDDINGS,
      metadatas: METADATAS,
    });
    try {
      await collection.get({
        where: {
          //@ts-ignore supposed to fail
          test: { $contains: "hello" },
        },
      });
    } catch (error: any) {
      expect(error).toBeDefined();
      expect(error.message).toMatchInlineSnapshot(
        `"Expected operator to be one of $gt, $gte, $lt, $lte, $ne, $eq, $in, $nin, but got $contains"`,
      );
    }
  });

  test("it should get embedding with matching documents", async () => {
    const collection = await client.createCollection({ name: "test" });
    await collection.add({
      ids: IDS,
      embeddings: EMBEDDINGS,
      metadatas: METADATAS,
      documents: DOCUMENTS,
    });
    const results2 = await collection.get({
      whereDocument: { $contains: "This is a test" },
    });
    expect(results2?.ids).toHaveLength(1);
    expect(["test1"]).toEqual(expect.arrayContaining(results2.ids));
  });

  test("it should get records not matching", async () => {
    const collection = await client.createCollection({ name: "test" });
    await collection.add({
      ids: IDS,
      embeddings: EMBEDDINGS,
      metadatas: METADATAS,
      documents: DOCUMENTS,
    });
    const results2 = await collection.get({
      whereDocument: { $not_contains: "This is another" },
    });
    expect(results2?.ids).toHaveLength(2);
    expect(["test1", "test3"]).toEqual(expect.arrayContaining(results2.ids));
  });

  test("test gt, lt, in a simple small way", async () => {
    const collection = await client.createCollection({ name: "test" });
    await collection.add({
      ids: IDS,
      embeddings: EMBEDDINGS,
      metadatas: METADATAS,
    });
    const items = await collection.get({
      where: { float_value: { $gt: -1.4 } },
    });
    expect(items.ids).toHaveLength(2);
    expect(["test2", "test3"]).toEqual(expect.arrayContaining(items.ids));
  });

  test("should error on non existing collection", async () => {
    const collection = await client.createCollection({ name: "test" });
    await client.deleteCollection({ name: "test" });
    await expect(async () => {
      await collection.get({ ids: IDS });
    }).rejects.toThrow();
  });

  test("it should throw an error if the collection does not exist", async () => {
    await expect(
      async () =>
        await client.getCollection({
          name: "test",
          embeddingFunction: new DefaultEmbeddingFunction(),
        }),
    ).rejects.toThrow(Error);
  });

  test("it should return results in row format", async () => {
    const collection = await client.createCollection({ name: "test" });
    await collection.add({
      ids: IDS,
      embeddings: EMBEDDINGS,
      metadatas: METADATAS,
      documents: DOCUMENTS,
    });

    const results = (
      await collection.get({ include: ["documents", "metadatas"] })
    ).rows();
    expect(results.length).toEqual(IDS.length);
    expect(results[0].document).toEqual("This is a test");
    expect(results[0].embedding).toBeUndefined();
    expect(results[0].metadata?.test).toEqual("test1");
  });
});
