# snnipp

This project is a minimal, private code snippet sharing application built with Go, MySQL, and Tailwind CSS. It's designed to be a personal tool for managing and sharing code snippets privately.

Here you can view it in action: [https://ssnipp.com](https://ssnipp.com)

## Features

- **Built with Go**: The application is powered by a Go web server.
- **MySQL Database**: Snippets are stored in a MySQL database.
- **Tailwind CSS**: The frontend is styled using Tailwind CSS.

## Requirements

- **.env File**: You need to create a `.env` file in the root directory to set up some default configurations.

## Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/benitolopez/ssnipp.git
   cd ssnipp
   ```

2. **Install dependencies**

   ```bash
   npm install
   ```

3. **Create the `.env` file**

   Create a `.env` file in the root directory of the project with the following content:

   ```env
   PORT=:4000
   DEBUG=false
   ALLOW_SIGNUP=true
   DB_USERNAME=your_db_username
   DB_PASSWORD=your_db_password
   DB_DATABASE=your_db_database
   DB_TEST_USERNAME=your_db_test_username
   DB_TEST_PASSWORD=your_db_test_password
   DB_TEST_DATABASE=your_db_test_database
   ```

   Replace `your_db_username`, `your_db_password`, and `your_db_name` (and test versions) with your actual MySQL credentials.

4. **Run the application**

   ```bash
   go run ./cmd/web
   ```

   or:

   ```bash
   npm run dev
   ```

   The application will start a web server, and you can access it via `http://localhost:4000`.

5. **Deploy the application**

   ```bash
   make build/web
   ```

## Note on Contributions

This project is maintained as my personal tool for private code snippet sharing. As such, I do not accept pull requests. However, feel free to fork the project and customize it for your own use.

## Acknowledgment

This project was inspired by the book [Let’s Go](https://lets-go.alexedwards.net/) by Alex Edwards. There are a few changes made in this implementation, but I totally recommend the book for anyone interested in learning Go and building web applications.
