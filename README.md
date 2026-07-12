# 🌐 AI Website Summarizer

A robust, full-stack web application designed to fetch website content and provide concise AI-driven summaries. This project prioritizes reliability and privacy by utilizing a dual-model architecture.

## 🚀 Key Features

* **HTML Scraping:** Cleans raw website data to extract readable text, ignoring scripts, styles, and navigation noise.
* **Dual-AI Architecture:**
    * **Cloud Mode:** Uses Google Gemini API for fast, high-quality summarization.
    * **Local Mode:** Uses Ollama (running `llama3`) for private, unlimited, and offline processing.
* **Intelligent Failover:** Automatically detects if the Gemini API is rate-limited or unavailable and seamlessly switches to your local Ollama instance with real-time UI updates.
* **Dynamic UI:** Real-time status notifications keep you informed when the system switches models or encounters errors.

## 🛠 Tech Stack

* **Backend:** Go (Golang)
* **Frontend:** HTML5, CSS3, Vanilla JavaScript
* **AI Engines:** Google Gemini API, Ollama (Llama3)
* **Infrastructure:** Local HTTP server

## 🧠 System Architecture



The application employs a "Traffic Controller" pattern where the frontend manages the failover logic, ensuring a smooth user experience even if a cloud provider goes down.

## ⚙️ Installation & Setup

### 1. Prerequisites
* [Go](https://go.dev/)
* [Ollama](https://ollama.com/)

### 2. Setup
1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd website-summarizer
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
**Scrape:** The Go backend fetches the URL and cleans the HTML tags using Regex to ensure only readable text is processed.

**Route:** The UI sends the URL and your preferred model to the backend.

**Summarize & Failover:**

* If Gemini is selected and succeeds, you get your result.

* If Gemini fails, the frontend dynamically alerts you ("Gemini failed, switching to local Ollama...") and triggers a secondary request to Ollama.

**Display:** The summary is injected into the clean UI with real-time status updates that automatically fade after the response is received.
