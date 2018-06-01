import React from 'react';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';

import { EditButton, RichTextField } from 'react-admin';

import InfoIcon from '@material-ui/icons/Info';

export const DashText = ({...props}) => {
    const { name, style, subtitle, txt } = props;
    props.basePath='board';
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
