import passport from 'passport';
import { Strategy as GoogleStrategy } from 'passport-google-oauth20';
import { Profile } from 'passport-google-oauth20';
import { UserModel } from '../models/user';
import { User } from '../types/user';
import dotenv from 'dotenv';

// Fix type augmentation to avoid recursive reference
declare global {
  namespace Express {
    interface User {
      id: number;
      google_id?: string;
      email: string;
      name: string;
      profile_picture?: string;
      created_at?: Date;
    }
  }
}

dotenv.config();

if (!process.env.GOOGLE_CLIENT_ID || !process.env.GOOGLE_CLIENT_SECRET) {
  throw new Error('Missing required Google OAuth credentials');
}

passport.use(new GoogleStrategy({
  clientID: process.env.GOOGLE_CLIENT_ID,
  clientSecret: process.env.GOOGLE_CLIENT_SECRET,
  callbackURL: "/auth/google/callback",
  proxy: true
}, async (
  accessToken: string,
  refreshToken: string,
  profile: Profile,
  done: (error: any, user?: any) => void
) => {
  try {
    let user = await UserModel.findByGoogleId(profile.id);

    if (!user) {
      await UserModel.create({
        google_id: profile.id,
        email: profile.emails?.[0]?.value || '',
        name: profile.displayName,
        profile_picture: profile.photos?.[0]?.value
      });
      
      user = await UserModel.findByGoogleId(profile.id);
    }

    return done(null, user);
  } catch (error) {
    return done(error as Error, undefined);
  }
}));

passport.serializeUser((user: Express.User, done) => {
  done(null, user.id);
});

passport.deserializeUser(async (id: number, done) => {
  try {
    const user = await UserModel.findById(id);
    if (!user) {
      return done(null, false);
    }
    return done(null,user as Express.User);
  } catch (error) {
    return done(error);
  }
});

export default passport;