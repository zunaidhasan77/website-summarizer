# 🌐 AI Website Assistant (RAG Enabled)

A robust, full-stack AI assistant that bridges the gap between static web content and intelligent, contextual conversation. This application allows users to perform broad summarizations or engage in deep, context-aware chats with any website.

## 🚀 Key Features

* **Dual-Mode AI:**
    * **Summarize Mode:** Instant, concise 3-paragraph summaries of any URL.
    * **Chat Mode (RAG):** Ask specific questions about a website. The system retrieves relevant context chunks and answers based on facts, not hallucinations.
* **Just-in-Time Ingestion:** No manual database setup required. The system automatically ingests and vectorizes the target website the moment you ask a question.
* **Intelligent Failover:** Automatically detects if the Gemini API is rate-limited or unavailable and seamlessly switches to your local Ollama instance with logging for transparency.
* **Vector Search:** Powered by **Qdrant**, enabling high-speed semantic search across website content.

## 🛠 Tech Stack

* **Backend:** Go (Golang)
* **Frontend:** HTML5, CSS3, Vanilla JavaScript
* **AI Engines:** Google Gemini API, Ollama (Llama3)
* **Vector Database:** Qdrant
* **Storage:** SQLite

## 🧠 System Architecture

The application utilizes a **RAG (Retrieval-Augmented Generation)** pipeline. When a user requests a chat:
1. **Ingest:** The URL is scraped and content is broken into searchable chunks.
2. **Vectorize:** Chunks are turned into embeddings and stored in Qdrant.
3. **Retrieve:** The system searches Qdrant for the most relevant context based on the user's question.
4. **Augment:** The final prompt sent to the LLM is injected with the retrieved context, ensuring the AI answers using only the provided website data.

## ⚙️ Installation & Setup

### 1. Prerequisites
* [Go](https://go.dev/)
* [Ollama](https://ollama.com/)
* [Qdrant](https://qdrant.tech/) (Docker: `docker run -p 6333:6333 qdrant/qdrant`)

### 2. Setup
1. **Clone the repository:**
   ```bash
   git clone https://github.com/zunaidhasan77/website-summarizer.git
   cd website-summarizer
   go mod tidy
2. **Download the local model:**
   ```bash
   ollama pull llama3
3. **Configure API Key:**
   Ensure your GEMINI_API_KEY is set in your environment variables.
4. **Run the Application**
   ```bash
   go run cmd/api/main.go
5. **Access the interface:**
   Open your browser to http://localhost:8080.

## How it Works
**Summarize Flow:** The backend scrapes the URL, truncates content to fit context limits, and requests a summary from the chosen AI model.

**Chat Flow** he backend performs a "Just-in-Time" ingestion (vectorizing the page), performs a semantic search to find relevant context, and sends a context-aware prompt to the LLM.

**Failover Logic:** If the primary model (Gemini) returns an error, the backend logs the failure and automatically retries the request using the local Ollama instance.


*Built with care for speed, reliability, and privacy.*
