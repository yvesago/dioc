import React from 'react';
import { useMediaQuery } from '@material-ui/core';
import { List, Datagrid, TextField, Edit, Create, SimpleForm, CreateButton,
    TextInput, required, EditButton, DateField, NumberInput, downloadCSV,
    RichTextField, Filter, TopToolbar, SimpleList, ExportButton
} from 'react-admin';
import jsonExport from 'jsonexport/dist';
import RichTextInput from 'ra-input-rich-text';
import { withStyles } from '@material-ui/core/styles';
import MyLeaflet from './Leaflet';


import ActionFlushIPButton from './ActionFlushIP';

const cardActionStyle = {
    zIndex: 2,
    display: 'inline-block',
    float: 'right',
};

const exporter = ips => {
    jsonExport(ips, {
        headers: ['name', 'host', 'count', 'p', 'asnnum']
    }, (err, csv) => {
        downloadCSV(csv, 'posts');
    });
};

const FlushIPActions = ({ resource, filters, displayedFilters, filterValues, basePath, showFilter, currentSort, exporter }) => (
    <TopToolbar style={cardActionStyle}>
        {filters && React.cloneElement(filters, { resource, showFilter, displayedFilters, filterValues, context: 'button' }) }
        <ActionFlushIPButton />
        <ExportButton 
            resource={resource}
            sort={currentSort}
            filter={filterValues}
            exporter={exporter}
            maxResults={10000}
        />
        <CreateButton basePath={basePath} />
    </TopToolbar>
);


const IPMap = ({ record }) => {
    return <div id="mapContainer"><MyLeaflet zoom={12} lat={record.lat} lng={record.lon} point={record.id} /></div>;
};

const IPFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Country" source="p" alwaysOn/>
        <TextInput label="Comment" source="comment" />
        <TextInput label="Name" source="name" />
        <TextInput label="Host" source="host" />
        <TextInput label="AS Num" source="asnnum" />
    </Filter>
);

const styles = {
    fieldASN: {
        display: 'inline-block', width: '100px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
    field: {
        display: 'inline-block', width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const IPList = withStyles(styles)(({ classes, ...props }) => {
    const isSmall = useMediaQuery(theme => theme.breakpoints.down('sm'));
    return (
        <List filters={<IPFilter />} perPage={30} actions={<FlushIPActions />} exporter={exporter} {...props}>
            {isSmall ? (
                <SimpleList
                    primaryText={record => `${record.name} -- ${record.host}`}
                    secondaryText={record => `${record.p} AS: ${record.asnname}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            ) : (
                <Datagrid>
                    <TextField source="name" />
                    <TextField source="host" />
                    <TextField source="count" />
                    <TextField label="Country" source="p" />
                    <TextField label="Region" source="r" />
                    <TextField label="City" source="c" />
                    <TextField label="AS Num" source="asnnum" />
                    <RichTextField source="asnname" className={classes.fieldASN} stripTags />
                    <RichTextField source="comment" className={classes.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
});


export const IPCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" validate={required()} />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const IPEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextField source="name" />
            <TextInput source="host" />
            <NumberInput source="count" />
            <TextInput label="AS Num" source="asnnum" />
            <TextInput label="AS Name" source="asnname" />
            <TextInput label="Country" source="p" />
            <TextInput label="Region" source="r" />
            <TextInput label="City" source="c" />
            <RichTextInput source="comment" />
            <DateField label="Created" source="created" showTime />
            <DateField label="Updated" source="updated" showTime />
            <NumberInput source="lat" />
            <NumberInput source="lon" />
            {<IPMap />}
        </SimpleForm>
    </Edit>
);
