import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, DisabledInput, required, EditButton, DateField, 
    RichTextField, SelectInput, LongTextInput, Filter
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

import { roles } from './MyConfig';


const cmd = [
    { label: 'Send Lines', id: 'SendLines', name: 'SendLines' },
    { label: 'Search Full File', id: 'FullSearch', name: 'FullSearch' },
    { label: 'STOP', id: 'STOP', name: 'STOP' },
];


const colored = WrappedComponent => props => props.record.status === 'OffLine' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'OnLine' ?
        <span style={{ color: 'green' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);

const AgentTitle = ({ record }) => {
    return <span>Agent {record ? `${record.ip}` : ''}</span>;
};


const AgentFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Ip" source="ip" alwaysOn />
        <TextInput label="Lines" source="lines" />
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} />
    </Filter>
);


export const AgentList = (props) => (
    <List filters={<AgentFilter />} {...props}>
        <Datagrid>
            <TextField source="ip" />
            <TextField source="filesurvey" />
            <ColoredTextField source="status" />
            <TextField source="role" />
            <TextField source="lines" />
            <RichTextField source="comment" stripTags />
            <TextField source="cmd" />
            <DateField label="updated" source="updated" showTime />
            <EditButton />
        </Datagrid>
    </List>
);


export const AgentCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="ip" validate={required} />
            <TextInput source="filesurvey" />
            <TextInput source="status" />
            <SelectInput source="role" choices={roles} allowEmpty />
            <TextInput source="lines" options={{ multiLine: true }}  />
            <LongTextInput source="comment" />
        </SimpleForm>
    </Create>
);

export const AgentEdit = (props) => (
    <Edit title={<AgentTitle />} {...props}>
        <SimpleForm>
            <TextField source="filesurvey" />
            <DateField label="created" source="created" showTime />
            <DateField label="updated" source="updated" showTime />
            <TextField source="lines"  style={{whiteSpace: 'pre-line'}}/>
            <SelectInput source="role" choices={roles} allowEmpty />
            <SelectInput source="cmd" choices={cmd} allowEmpty optionText="label" />
            <RichTextInput source="comment" />
            <DisabledInput source="crca" />
        </SimpleForm>
    </Edit>
);
