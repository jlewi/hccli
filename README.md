# hccli
CLI for honeycomb and observability queries

## Usage instructs

1. Set the endpoint of your AI that converts natural language into HoneycombQueries
   ```bash
    hccli config set aiEndpoint=http://localhost:5000
    ```
   
    * This should be an instance of this [cog server](https://github.com/hamelsmu/replicate-examples/tree/79ec0e71b120dc1bcf6c3c7b26f9331e9e734f2a/mistral-vllm-awq)

1. Set the path to a file containing your honeycomb API key
   
    ```bash
    hccli  config set honeycombApiKeyFile=~/.honeycomb_api_key
    ```
   
## Limitations

Unfortunately the Honeycomb API only lets you fetch query data if your on the enterprise plan.
See [QueryData API Docs](https://docs.honeycomb.io/api/tag/Query-Data)