# EquiTrack

EquiTrack is a stock tracking platform designed to help users manage their investments efficiently. With EquiTrack, users can add, remove, edit, and view their investment records, as well as monitor real-time stock returns for their portfolio.

## Features

- **Investment Management**: Add, Edit, and Delete investment records seamlessly.
- **Real-Time Tracking**: View real-time stock returns and the current value of your portfolio.
- **Expandable Details**: Click on a stock to view detailed individual investments.
- **Interactive Dashboard**: A clean and intuitive UI to provide summaries of total investments, current value, and profit & loss.

## Tech Stack

- **Frontend**: React, Ant Design (AntD) for a responsive and visually appealing user interface.
- **Backend**: Golang for building a robust and scalable API.
- **Styling**: Custom CSS with responsive design to ensure compatibility across devices.

## Installation and Setup

Follow these steps to set up the project locally:

### Prerequisites

- Node.js (v14+)
- Golang (v1.19+)
- MySQL (or any compatible database)

### Backend Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/equitrack.git
   cd equitrack/backend
   ```

2. Configure the database:

   - Update the database credentials in the `config` directory.

3. Install dependencies and run the server:

   ```bash
   go mod tidy
   go run main.go
   ```

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd ../frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```

The application will now be accessible at `http://localhost:3000`.

## API Endpoints

### Auth API

- **POST /login**
  - Check the login credentials of the user.
  - Query Params: `username`,`password`
- **POST /register**
  - Adds new user to database.
  - Query Params: `username`,`email`,`password`

### Investments API

- **GET /investments**
  - Fetch all investments for a user.
  - Query Params: `userId`
- **GET /individualInvestments**
  - Fetch detailed individual investments for a specific stock.
  - Query Params: `userId`, `instrument`
- **POST /investments**
  - Fetch all investments for a user.
  - Query Params: `userId`,`instrument`,`quantity`,`avg_price`
- **DELETE /investments**
  - Remove an investment.
  - Query Params: `id`

## Screenshots

### Login Page
![image](https://github.com/user-attachments/assets/df0e7201-e42b-4d2d-92bc-316e334a42a8)

### Dashboard
![image](https://github.com/user-attachments/assets/8725e85a-dd21-471d-a37c-7467a17fb143)




## Contributing

We welcome contributions to EquiTrack! Here's how you can get involved:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-name`.
3. Commit your changes: `git commit -m 'Add some feature'`.
4. Push to the branch: `git push origin feature-name`.
5. Open a pull request.


## Contact

For questions or support, please contact [danish.eqbal125@gmail.com](mailto\:danish.eqbal125@gmail.com).

