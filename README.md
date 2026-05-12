## Setup

1. **Install Go** (version 1.21 or later)

2. **Get your OpenRouter API key**:
   - Visit [OpenRouter.AI](https://openrouter.ai/)
   - Sign up and get your API key

3. **Set environment variable**:
   ```bash
   export OPENROUTER_API_KEY="your-api-key-here"
   ```

4. **Install dependencies**:
   ```bash
   go mod tidy
   ```

5. **Run the interactive chat**:
   ```bash
   go run main.go
   ```