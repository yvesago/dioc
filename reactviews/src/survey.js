import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, DisabledInput, required, EditButton, DateField, 
    RichTextField, SelectInput 
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

import { roles } from './MyConfig';

const levels = [
    { name: 'Critic', id: 'critic' },
    { name: 'Warn', id: 'warn' },
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