# Demo Simple Orc Agent

The personality of the Orc is based on the character sheet provided in `character_sheet.md`.
The `instructions.md` guides the Orc to answer questions based on this character sheet.

## Model

```bash
docker model pull ai/qwen2.5:latest
```
> If you use Docker Compose, this will be pulled automatically.

## Start the application

**With Docker Compose**:
```bash
docker compose up --build -d
docker attach $(docker compose ps -q orc-agent)
```
> `docker compose down` to stop the application.


**From a container**:
```bash
MODEL_RUNNER_BASE_URL=http://model-runner.docker.internal/engines/llama.cpp/v1 \
MODEL_RUNNER_CHAT_MODEL=ai/qwen2.5:latest \
go run main.go
```


**From a local machine**:
```bash
MODEL_RUNNER_BASE_URL=http://localhost:12434/engines/llama.cpp/v1 \
MODEL_RUNNER_CHAT_MODEL=ai/qwen2.5:latest \
go run main.go
```

## Talk to the Orc

- What is your name?
- What is your occupation?
- What is your favorite food?  
- What is your background story?
- What is your main quote?

## Information

### Budgie

This program has been created thanks to the **budgie** library: [https://github.com/budgies-nest/budgie](https://github.com/budgies-nest/budgie)

### Devcontainer
This project is designed to run in a [Devcontainer](https://code.visualstudio.com/docs/devcontainers/containers) environment.
You can use the provided `.devcontainer` configuration in the `snippets` directory to set up the environment in Visual Studio Code or any compatible IDE.
