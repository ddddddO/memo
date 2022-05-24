import * as React from 'react';
import { createRoot } from 'react-dom/client';
import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider } from '@mui/material/styles';
import App from './App';
import theme from './theme';
import Memos from './Memos';
import SignIn from './SignIn';
import { BrowserRouter, Routes, Route/**, Link**/ } from "react-router-dom";

const rootElement = document.getElementById('root');
const root = createRoot(rootElement!);

root.render(
    <ThemeProvider theme={theme}>
      <BrowserRouter>
        <App />
        <Routes>
          <Route path="" element={<SignIn />} />
          <Route path="/signin" element={<SignIn />} />
          <Route path="/memos" element={<Memos />} />
        </Routes>
      </BrowserRouter>
      {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
      <CssBaseline />
    </ThemeProvider>
);