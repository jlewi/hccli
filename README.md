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

## Visualizing Honeycomb Queries

You can use [Honeycomb's Query Sharing Feature](https://docs.honeycomb.io/investigate/collaborate/share-query/)
to generate URLs that have the query directly in the browser which you can then just open. The CLI
allows you to open these URLs directly from the command line and take a snapshot using
[chromedp](https://github.com/chromedp/chromedp).

1. Start Chrome with remote debugging enabled

     ```bash
     chrome --remote-debugging-port=9222
     ```

1. Open up your browser and login to Honeycomb

1. Use the CLI to open the URL in the browser and take a snapshot

    ```bash
    hccli --query-file=model_query.json --dataset=service --base-url=https://ui.honeycomb.io/autobuilder/environments/prod/datasets/production --out-file=/tmp/screenshot.png
    ```

## Limitations

Unfortunately the Honeycomb API only lets you fetch query data if your on the enterprise plan.
See [QueryData API Docs](https://docs.honeycomb.io/api/tag/Query-Data)