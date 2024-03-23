import React from 'react';
import { fetchUtils, Admin, Resource } from 'react-admin';
import { Route } from 'react-router-dom';
//import './App.css';

import AlertIcon from '@mui/icons-material/Assessment';
import SurveyIcon from '@mui/icons-material/WifiTethering';
import AgentIcon from '@mui/icons-material/Wifi';
import ExtractIcon from '@mui/icons-material/Transform';
import IpIcon from '@mui/icons-material/Fingerprint';
import MapIcon from '@mui/icons-material/Public';

import { createTheme } from '@mui/material/styles';
import myTheme from './myTheme';

import { Dashboard } from './Dashboard';
import { AlertList, AlertEdit } from './alert';
import { SurveyList, SurveyEdit, SurveyCreate } from './survey';
import { ExtractList, ExtractEdit, ExtractCreate } from './extract';
import { IPList, IPEdit, IPCreate } from './ip';
import { AgentList, AgentEdit, AgentCreate } from './agent';
import { BoardEdit } from './boarddocs';
import Mapview from './Mapview';

import { MyConfig } from './MyConfig';
import mySimpleRest from './myJsonRestNew';

import authLoginPage from './authLoginPage';
import authClient from './authClient';

const httpClient = (url, options = {}) => {
    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }
    // add your own headers here
    //options.headers.set('X-MyToken', MyConfig.API_KEY );
    const token = localStorage.getItem('gotoken');
    options.headers.set('Authorization', `Bearer ${token}`);
    return fetchUtils.fetchJson(url, options);
};


const restClient = mySimpleRest( MyConfig.API_URL  + '/admin/api/v1', httpClient);


const App = () => (
    <Admin disableTelemetry title='Distributed IOC manager' theme={ createTheme(myTheme) }
        loginPage={authLoginPage} authProvider={authClient(MyConfig.AUTH_URL)}
        dashboard={Dashboard} 
        dataProvider={restClient}>
        <Resource name="alertes" list={AlertList} edit={AlertEdit} icon={AlertIcon} />
        <Resource name="surveys" list={SurveyList}  edit={SurveyEdit} create={SurveyCreate} icon={SurveyIcon} />
        <Resource name="agents" list={AgentList}  edit={AgentEdit} create={AgentCreate} icon={AgentIcon} />
        <Resource name="extracts" list={ExtractList}  edit={ExtractEdit} create={ExtractCreate} icon={ExtractIcon} />
        <Resource name="ips" options={{ label: '‣ Extract > IPs'}} list={IPList}  edit={IPEdit} create={IPCreate} icon={IpIcon} />
        <Resource name='Map' options={{ label: '‣ Extract > Map'}} list={Mapview} icon={MapIcon} />
        <Resource name="board" edit={BoardEdit} />
    </Admin>
);

export default App;

