AI Website Summarizer
A robust, full-stack web application that fetches website content and provides concise summaries using AI. The application is designed for reliability, featuring model-switching between cloud-based AI and local models to ensure you never face downtime.

🚀 Key Features
HTML Scraping: Cleans raw website data to extract readable text, ignoring scripts, styles, and navigation junk.

Dual-AI Architecture: * Cloud Mode: Uses Gemini API for fast, high-quality summarization.

Local Mode: Uses Ollama (running llama3) for private, unlimited, and offline summarization.

Intelligent Failover: Automatically detects if the Gemini API is rate-limited or unavailable and seamlessly switches to your local Ollama instance.

Dynamic UI: Real-time status notifications keep the user informed when the system switches models.

🛠 Tech Stack
Backend: Go (Golang)

Frontend: HTML5, CSS3, Vanilla JavaScript

AI Engines: Google Gemini API, Ollama (Llama3)

Scraping: Standard Go net/http with Regex text cleaning

⚙️ Installation & Setup
1. Prerequisites
Go installed on your system.

Ollama installed and running locally.

A valid Gemini API Key.

2. Setup
Clone the repository:

Bash
git clone <your-repo-url>
cd website-summarizer
Pull the required local model:

Bash
ollama pull llama3
Set your Gemini API key (ensure your environment variable is set in your main.go or config):

Bash
export GEMINI_API_KEY="your_api_key_here"
3. Running the App
Start the backend:

Bash
go run cmd/api/main.go
Open your browser and navigate to http://localhost:8080.

🧠 Architecture Overview
The application follows a "Traffic Controller" pattern where the frontend manages the failover logic, ensuring the user experience remains smooth even if one service provider goes down.

📝 License
This project is open-source and intended for personal learning and productivity enhancement.

This README covers everything you've built! Does this look like a good summary of your project, or would you like to add anything else?