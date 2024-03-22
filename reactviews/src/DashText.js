import React from 'react';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import CardContent from '@mui/material/CardContent';

import { EditButton, RichTextField } from 'react-admin';

import InfoIcon from '@mui/icons-material/Info';

export const DashText = ({...props}) => {
    const { name, style, subtitle, txt } = props;
    props.resource='board';
    return (
        <Card style={style.card}>
            <InfoIcon style={style.icon} />
            <CardHeader
                title={name}
                subheader={subtitle}
            />
            <CardContent>
                <RichTextField source={txt} {...props} />
                <EditButton { ...props} />
            </CardContent>
        </Card>
    );
};
