import { Router, Request, Response } from 'express';
import passport from 'passport';

const authRoute = Router();

// Authentication check middleware
const isAuthenticated = (req: Request, res: Response, next: Function) => {
  if (req.isAuthenticated()) {
    return next();
  }
  res.status(401).json({ error: 'Not authenticated' });
};

// Google OAuth routes
authRoute.get('/google', 
  passport.authenticate('google', { 
    scope: ['profile', 'email'],
    prompt: 'select_account'
  })
);

authRoute.get('/google/callback', 
  passport.authenticate('google', { 
    failureRedirect: '/auth/login-failed',
    successRedirect: process.env.CLIENT_URL || 'http://localhost:5173/wallet'
  })
);

// Additional auth routes
authRoute.get('/login-failed', (req: Request, res: Response) => {
  res.status(401).json({
    success: false,
    message: 'Authentication failed'
  });
});

authRoute.get('/logout', (req: Request, res: Response) => {
  req.logout(() => {
    res.redirect(process.env.CLIENT_URL || 'http://localhost:5173');
  });
});

// Check auth status
authRoute.get('/status', (req: Request, res: Response) => {
  res.json({
    isAuthenticated: req.isAuthenticated(),
    user: req.user
  });
});

export default authRoute;