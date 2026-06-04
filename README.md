Neon Chat 🚀

A real-time chat application built using Go (Golang) and WebSockets, with a simple and clean frontend using HTML, CSS, and JavaScript.

💡 About the Project

Neon Chat is a lightweight real-time messaging app where users can join a chat room, set their own username, and communicate instantly with others.

It focuses on understanding how real-time communication works using WebSockets.

⚙️ Features
Real-time messaging
Custom usernames
Online users list
Join / leave notifications
Typing indicator
Clean neon-style UI

🛠️ Tech Stack
Go (Golang)
WebSockets (gorilla/websocket)
HTML
CSS
JavaScript

📂 Project Structure
chatapp/
│
├── client/        # Frontend files
├── server/        # Go backend (WebSocket server)
├── README.md
├── run_windows.bat
└── run_linux_mac.sh


🚀 How to Run
1. Clone the project
git clone https://github.com/your-username/neon-chat.git

2. Run Backend
cd server
go run .

Server will start at:

http://localhost:8080


3. Run Frontend

Open client/index.html in browser
OR use Live Server in VS Code.


🧠 Learning Purpose

This project was built to understand:
WebSocket communication
Real-time data flow
Client-server interaction
Basic full-stack architecture using Go


📸 UI

Simple neon-style interface with smooth real-time chat experience.


👨‍💻 Author

Built with ❤️ for learning real-time systems and improving backend skills.

📌 Note

This is a learning project and not intended for production use.

