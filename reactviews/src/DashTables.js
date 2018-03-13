import React from 'react';
import { Card, CardTitle } from 'material-ui/Card';
//import { List, ListItem } from 'material-ui/List';
import TableIcon from 'material-ui/svg-icons/action/assessment';

import ChartDoghnut from './ChartDoghnut.js';

const styles = {
    card: { borderLeft: 'solid 4px #ff9800', flex: 1, marginLeft: '1em' },
    icon: { float: 'right', width: 64, height: 64, padding: 16, color: '#ff9800' },
    dognut: { float: 'left', display: 'inline-block', marginLeft: 15 },
};

export default ({ nbagents = [], nbsurveys = [], nbalerts = [] }) => (
    <Card style={styles.card}>
        <TableIcon style={styles.icon} />
        <CardTitle title="Tables" />
        <div style={styles.dognut} >
            <ChartDoghnut title="Agents" data={nbagents} />
        </div>
        <div style={styles.dognut} >
            <ChartDoghnut title="Surveys" data={nbsurveys} />
        </div>
        <div style={styles.dognut} >
            <ChartDoghnut title="Alerts" data={nbalerts} />
        </div>
    </Card>
);
