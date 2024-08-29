> NOTE: This project is currently unstable and under active development

# Pinecone Mock Server

Lightweight server that imitates the [Pinecone REST API](https://docs.pinecone.io/reference/api/introduction)

Intended to be used for integration tests where it's not suitable to call the Live Pinecone API.

No dependencies, just pure Go

## Supports

| Operation | Pinecone Endpoint | Supported |
| --------- | ----------------- | --------- |
| Create Index | `POST /indexes`  | **Yes** |
| List Indexes | `GET /indexes`   | **Yes** |
| Describe Index | `GET /indexes/{index_name}` | **Yes** |
| Delete Index | `DELETE /indexes/{index_name}` | No |
| Configure Index | `PATCH /indexes/{index_name}` | No |
| Get Index Stats| `GET /describe_index_stats` | **Yes** |
| Upsert Vectors | `POST /vectors/upsert` | **Yes** |
| Query Vectors | `POST /query` | **Yes** |
| Fetch Vectors | `GET /vectors/fetch` | **Yes** |
| Update Vectors | `POST /vectors/update` | **Yes** |
| Delete Vectors | `POST /vectors/delete` | **Yes** |
| List Vectors | `GET /vectors/list` | **Yes** |
| List Collection | `GET /collections` | No |
| Create Collection | `POST /collections` | No |
| Describe Collection | `GET /collections/{collection_name}` | No |
| Delete Collection | `GET /collections/{collection_name}` | No |


# Usage

```
$ go build -o mock main.go

$ ./mock
Server listening on port :8080
```

Now you can point Pinecone to the server running on `:8080`

```python
from pinecone import Pinecone

client = Pinecone(host="localhost:8080", api_key="test")

client.create_index("test_index", dimension=3, spec={})

index = client.Index("test_index")

index.upsert([('id1', [1.0, 2.0, 3.0], {'key': 'value'}), ('id2', [1.0, 2.0, 3.0])], namespace="test")

index.update(id="id1", values=[4.0, 5.0, 6.0], set_metadata={"genre": "comedy"}, namespace="test")

next(index.list())
['id1', 'id2']
```

You can also use the REST API:

```
≫ curl -XPOST "http://localhost:8080/query" -d'{"topK": 2}' | jq
{
  "matches": [
    {
      "id": "id1",
      "values": [
        4,
        5,
        6
      ],
      "sparseValues": {},
      "metadata": {
        "key": "value",
        "genre": "comedy"
      }
    },
    {
      "id": "id2",
      "values": [
        1,
        2,
        3
      ],
      "sparseValues": {}
    }
  ],
  "namespace": "",
  "usage": {
    "read_units": 2
  }
}
```
