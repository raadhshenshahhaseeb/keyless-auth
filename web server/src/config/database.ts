import mysql from 'mysql2/promise';
import dotenv from 'dotenv';
import path from 'path';


dotenv.config();
export const pool = mysql.createPool({
  host: 'localhost', 
  port:3306,
  user: 'root',
  password: process.env.DB_PASSWORD,
  database: 'zklogin_db',
  waitForConnections: true,
  connectionLimit: 10,
  queueLimit: 0,
  authPlugins: {
    mysql_native_password: () => () => {
      return Buffer.from(process.env.DB_PASSWORD || '')
    }
  }
});

// Test connection function
async function testConnection() {
  try {
    const connection = await pool.getConnection();
    console.log('Database configuration:', {
      host: connection.config.host,
      user: connection.config.user,
      database: connection.config.database
    });
    
    const [rows] = await connection.query('SELECT 1 as result');
    console.log('✅ Database connected successfully', rows);
    
    connection.release();
    return true;
  } catch (error: any) {
    console.error('Connection error details:', {
      message: error.message,
      code: error.code,
      errno: error.errno,
      sqlState: error.sqlState,
      host: process.env.DB_HOST,
      user: process.env.DB_USER,
      database: process.env.DB_NAME
    });
    return false;
  }
}

// Execute test connection
testConnection();

export default pool;