import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, DisabledInput, required, EditButton, DateField, 
    RichTextField, SelectInput 
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

const levels = [
    { name: 'Critic', id: 'critic' },
    { name: 'Warn', id: 'warn' },
];
    //{ label: 'None', value: '' }

const roles = [ // XXX fix your categorys
    { name: 'Squid', id: 'squid' },
    { name: 'Radius', id: 'radius' },
    { name: 'Web', id: 'web' },
    { name: 'DNS', id: 'dns' },
    { name: 'Honeypot', id: 'honeypot' },
    { name: 'NetFlow', id: 'netflow' },
    { name: 'Auth', id: 'auth' },
    { name: 'SMTP', id: 'smtp' },
    { name: 'Mail', id: 'mail' },
    { name: 'Test', id: 'test' },
];


const colored = WrappedComponent => props => props.record.level === 'critic' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'warn' ? 
        <span style={{ color: 'orange' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);


export const SurveyList = (props) => (
    <List {...props}>
        <Datagrid>
            <TextField source="search" />
            <TextField source="role" />
            <ColoredTextField source="level" />
            <RichTextField source="comment" stripTags />
            <DateField label="updated" source="updated" showTime />
            <EditButton />
        </Datagrid>
    </List>
);


export const SurveyCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required} />
            <SelectInput source="level" choices={levels} allowEmpty />
            <SelectInput source="role" choices={roles} allowEmpty />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);
//<DateInput label="Publication date" source="published_at" defaultValue={new Date()} />

export const SurveyEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required} />
            <DateField label="Created" source="created" showTime />
            <DateField label="Updated" source="updated" showTime />
            <SelectInput source="role" choices={roles} allowEmpty />
            <SelectInput source="level" choices={levels} allowEmpty />
            <RichTextInput source="comment" />
            <DisabledInput source="crcs" validate={required} />
        </SimpleForm>
    </Edit>
);
