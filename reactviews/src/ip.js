import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField, NumberInput,
    RichTextField, Filter, Responsive, SimpleList,
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

import MyLeaflet from './Leaflet';

/*const colored = WrappedComponent => props => props.record.actions === 'Delete' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'warn' ? 
        <span style={{ color: 'orange' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);
*/

const IPMap = ({ record }) => {
    return <div id="mapContainer"><MyLeaflet zoom={12} lat={record.lat} lng={record.lon} point={record.id} /></div>;
};

const IPFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Country" source="p" alwaysOn/>
        <TextInput label="Comment" source="comment" />
        <TextInput label="AS Num" source="asnnum" />
    </Filter>
);


export const IPList = (props) => (
    <List filters={<IPFilter />} perPage={30} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => `${record.name} -- ${record.host}`}
                    secondaryText={record => `${record.p} AS: ${record.asname}`}
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
                    <TextField label="AS Num" source="asnum" />
                    <RichTextField source="asname" stripTags elStyle={{width: '100px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}} />
                    <RichTextField source="comment" stripTags elStyle={{width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}} />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
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