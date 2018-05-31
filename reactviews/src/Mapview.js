import React from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { ViewTitle } from 'react-admin';
import MyLeaflet from './Leaflet';

export default () => (
    <Card>
        <ViewTitle title="IP map" />
        <CardContent>From active Extracts</CardContent>
        <div id="mapContainer">
            <MyLeaflet zoom={4} />
        </div>

    </Card>
);

