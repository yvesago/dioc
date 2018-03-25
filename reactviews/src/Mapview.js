import React from 'react';
import { Card, CardText } from 'material-ui/Card';
import { ViewTitle } from 'admin-on-rest/lib/mui';
import MyLeaflet from './Leaflet';

export default () => (
    <Card>
        <ViewTitle title="IP map" />
        <CardText>From active Extracts</CardText>
        <div id="mapContainer">
            <MyLeaflet zoom={4} />
        </div>

    </Card>
);

