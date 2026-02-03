import React, { createContext, useContext, useState, useEffect } from 'react';
import enUS from '../locales/en-US';
import ptBR from '../locales/pt-BR';

const PreferencesContext = createContext();

const translations = {
    'en-US': enUS,
    'pt-BR': ptBR
};

export const PreferencesProvider = ({ children }) => {
    // Theme State
    const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'light');

    // Language State
    const [language, setLanguage] = useState(() => localStorage.getItem('language') || 'en-US');

    useEffect(() => {
        localStorage.setItem('theme', theme);
        document.documentElement.setAttribute('data-theme', theme);
    }, [theme]);

    useEffect(() => {
        localStorage.setItem('language', language);
    }, [language]);

    const toggleTheme = () => {
        setTheme(prev => prev === 'light' ? 'dark' : 'light');
    };

    // Translation Hook
    const t = (key) => {
        const dict = translations[language] || translations['en-US'];
        return dict[key] || key;
    };

    return (
        <PreferencesContext.Provider value={{ theme, toggleTheme, language, setLanguage, t }}>
            {children}
        </PreferencesContext.Provider>
    );
};

export const usePreferences = () => useContext(PreferencesContext);
