import React from 'react';
import { Route } from 'react-router-dom';
import  authCallBack  from './authCallback';

export default [
    <Route exact path="/callback" component={authCallBack} noLayout />,
];
