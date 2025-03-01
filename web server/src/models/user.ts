import { pool } from '../config/database';
import { User } from '../types/user';

export const UserModel = {
  findByGoogleId: async (googleId: string): Promise<User | undefined> => {
    const [rows] = await pool.execute(
      'SELECT * FROM users WHERE google_id = ?',
      [googleId]
    );
    return (rows as User[])[0];
  },

  create: async (user: User) => {
    // Update the SQL query to match all columns
    const [result] = await pool.execute(
      'INSERT INTO users (google_id, email, name, profile_picture) VALUES (?, ?, ?, ?)',
      [
        user.google_id,
        user.email,
        user.name,
        user.profile_picture
      ]
    );
    return result;
  },

  findById: async (id: number): Promise<User | undefined> => {
    const [rows] = await pool.execute(
      'SELECT * FROM users WHERE id = ?',
      [id]
    );
    return (rows as User[])[0];
  }
};