import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm, CreateButton,
    TextInput, required, EditButton, DateField, NumberInput, DeleteButton,
    RichTextField, Filter, Responsive, SimpleList, RefreshButton
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

import MyLeaflet from './Leaflet';

import { CardActions } from 'material-ui/Card';

import ActionFlushIPButton from './ActionFlushIP';

const cardActionStyle = {
    zIndex: 2,
    display: 'inline-block',
    float: 'right',
};

const FlushIPActions = ({ resource, filters, displayedFilters, filterValues, basePath, showFilter }) => (
    <CardActions style={cardActionStyle}>
        {filters && React.cloneElement(filters, { resource, showFilter, displayedFilters, filterValues, context: 'button' }) }
        <ActionFlushIPButton />
        <CreateButton basePath={basePath} />
        <RefreshButton />
    </CardActions>
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


export const IPList = (props) => (
    <List filters={<IPFilter />} perPage={30} actions={<FlushIPActions />} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => `${record.name} -- ${record.host}`}
                    secondaryText={record => `${record.p} AS: ${record.asnname}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            }
            medium={
                <Datagrid>
                    <TextField source="name" />
                    <TextField source="host" />
                    <TextField source="count" />
                    <TextField label="Country" source="p" />
                    <TextField label="Region" source="r" />
                    <TextField label="City" source="c" />
                    <TextField label="AS Num" source="asnnum" />
                    <RichTextField source="asnname" stripTags elStyle={{width: '100px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}} />
                    <RichTextField source="comment" stripTags elStyle={{width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}} />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                    <DeleteButton />
                </Datagrid>
            }
        />
    </List>
);


export const IPCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" validate={required} />
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
