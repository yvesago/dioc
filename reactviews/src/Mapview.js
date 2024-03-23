import React from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import { Title } from 'react-admin';
import MyLeaflet from './Leaflet';



export default () => (
    <Card>
        <Title title="IP map" />
        <CardContent>From active Extracts</CardContent>
        <div>
            <MyLeaflet zoom={4} />
        </div>

    </Card>
);

