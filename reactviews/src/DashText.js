import React from 'react';
import { Card, CardHeader, CardText } from 'material-ui/Card';
import { EditButton, RichTextField } from 'admin-on-rest';

import InfoIcon from 'material-ui/svg-icons/action/info';

export const DashText = ({...props}) => {
    const { name, style, subtitle, txt } = props;
    props.basePath='board';
    return (
        <Card style={style.card}>
            <InfoIcon style={style.icon} />
            <CardHeader
                title={name}
                subtitle={subtitle}
            />
            <CardText>  
                <RichTextField source={txt} {...props} />
                <EditButton { ...props} />
            </CardText>
        </Card>
    );
};
