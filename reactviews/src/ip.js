import React from 'react';
import { Card, useMediaQuery } from '@mui/material';
import { List, Datagrid, TextField, Edit, Create, SimpleForm, CreateButton,
    TextInput, required, EditButton, DateField, NumberInput, downloadCSV,
    RichTextField, TopToolbar, SimpleList, ExportButton, FilterButton,
    useRecordContext, Labeled
} from 'react-admin';
import jsonExport from 'jsonexport/dist';
//import { RichTextInput } from 'ra-input-rich-text';
const RichTextInput = React.lazy(() =>
    import('ra-input-rich-text').then(module => ({
        default: module.RichTextInput,
    }))
);
import MyLeaflet from './Leaflet';


import ActionFlushIPButton from './ActionFlushIP';

const exporter = ips => {
    jsonExport(ips, {
        headers: ['name', 'host', 'count', 'p', 'asnnum']
    }, (err, csv) => {
        downloadCSV(csv, 'posts');
    });
};

const FlushIPActions = () => (
    <TopToolbar>
        <FilterButton />
        <ActionFlushIPButton />
        <ExportButton />
        <CreateButton />
    </TopToolbar>
);


const IPMap = () => {
    const record = useRecordContext();
    return <Card sx={{width: '100%'}}> <MyLeaflet zoom={12} lat={record.lat} lng={record.lon} point={record.id} name={record.name} /></Card>;
};


const postFilters = [
    <TextInput label="Country" source="p" alwaysOn />,
    <TextInput label="Comment" source="comment" />,
    <TextInput label="Name" source="name" />,
    <TextInput label="Host" source="host" />,
    <TextInput label="AS Num" source="asnnum" />,
];

const styles = {
    fieldASN: {
        display: 'inline-block', width: '100px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'},
    field: {
        display: 'inline-block', width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const IPList = ({ props }) => {
    const isSmall = useMediaQuery(theme => theme.breakpoints.down('sm'));
    return (
        <List filters={postFilters} perPage={50} actions={<FlushIPActions />} exporter={exporter} {...props}>
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
                    <RichTextField source="asnname" sx={styles.fieldASN} stripTags />
                    <RichTextField source="comment" sx={styles.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
};


export const IPCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" validate={required()} />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const IPEdit = () => (
    <Edit>
        <SimpleForm>
            <Labeled label="Name">
                <TextField source="name" />
            </Labeled>
            <TextInput source="host" />
            <NumberInput source="count" />
            <TextInput label="AS Num" source="asnnum" />
            <TextInput label="AS Name" source="asnname" />
            <TextInput label="Country" source="p" />
            <TextInput label="Region" source="r" />
            <TextInput label="City" source="c" />
            <RichTextInput source="comment" />
            <Labeled label="Created">
                <DateField source="created" showTime />
            </Labeled>
            <Labeled label="Updated">
                <DateField source="updated" showTime />
            </Labeled>
            <NumberInput source="lat" />
            <NumberInput source="lon" />
            <IPMap />
        </SimpleForm>
    </Edit>
);
