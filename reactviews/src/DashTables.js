import React from 'react';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import TableIcon from '@mui/icons-material/Assessment';

import ChartDoghnut from './ChartDoghnut.js';

const styles = {
    card: { borderLeft: 'solid 4px #ff9800', flex: 1, marginLeft: '1em' },
    icon: { float: 'right', width: 64, height: 64, padding: 16, color: '#ff9800' },
    dognut: { float: 'left', display: 'inline-block', marginLeft: 15 },
};

export const DashTables = ({ nbagents = [], nbsurveys = [], nbalerts = [], subtitle }) => (
    <Card style={styles.card}>
        <TableIcon style={styles.icon} />
        <CardHeader title="Tables" subheader={subtitle} />
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
