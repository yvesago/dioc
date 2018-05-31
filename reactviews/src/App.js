import React from 'react';
import { fetchUtils, Admin, Resource } from 'react-admin';
//import './App.css';

import AlertIcon from '@material-ui/icons/Assessment';
import SurveyIcon from '@material-ui/icons/WifiTethering';
import AgentIcon from '@material-ui/icons/Wifi';
import SettingsIcon from '@material-ui/icons/Settings';
import MapIcon from '@material-ui/icons/Public';

import { createMuiTheme } from '@material-ui/core/styles';
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
import mySimpleRest from './myJsonRestServer';

import authLoginPage from './authLoginPage';
import authClient from './authClient';
import customRoutes from './customRoutes';

const httpClient = (url, options = {}) => {
    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }
    // add your own headers here
    //options.headers.set('X-MyToken', MyConfig.API_KEY );
    const token = localStorage.getItem('token');
    options.headers.set('Authorization', `Bearer ${token}`);
    return fetchUtils.fetchJson(url, options);
};


const restClient = mySimpleRest( MyConfig.API_URL  + '/admin/api/v1', httpClient);


const App = () => (
    <Admin title='Distributed IOC manager' theme={ createMuiTheme(myTheme) }
        loginPage={authLoginPage} authProvider={authClient(MyConfig.AUTH_URL)}
        customRoutes={customRoutes} dashboard={Dashboard} 
        dataProvider={restClient}>
        <Resource name="alertes" list={AlertList} edit={AlertEdit} icon={AlertIcon} />
        <Resource name="surveys" list={SurveyList}  edit={SurveyEdit} create={SurveyCreate} icon={SurveyIcon} />
        <Resource name="agents" list={AgentList}  edit={AgentEdit} create={AgentCreate} icon={AgentIcon} />
        <Resource name="extracts" list={ExtractList}  edit={ExtractEdit} create={ExtractCreate} icon={SettingsIcon} />
        <Resource name="ips" options={{ label: '‣ Extract > IPs'}} list={IPList}  edit={IPEdit} create={IPCreate} icon={AgentIcon} />
        <Resource name='Map' options={{ label: '‣ Extract > Map'}} list={Mapview} icon={MapIcon} />
        <Resource name="board" edit={BoardEdit} />
    </Admin>
);

export default App;

