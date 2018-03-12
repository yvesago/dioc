import React from 'react';
import { List, Datagrid, TextField, DateField, EditButton,
    Edit, SimpleForm, Filter, TextInput, SelectInput, DeleteButton
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

import { roles } from './MyConfig';

const colored = WrappedComponent => props => props.record.level === 'critic' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'warn' ?
        <span style={{ color: 'orange' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);

const AlertFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Lines" source="lines" />
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} allowEmpty alwaysOn />
    </Filter>
);


export const AlertList = (props) => (
    <List filters={<AlertFilter />} {...props}>
        <Datagrid>
            <TextField source="search" />
            <TextField source="role" />
            <TextField source="ip" />
            <TextField source="filesurvey" />
            <ColoredTextField source="level" />
            <TextField source="line" />
            <TextField source="comment" />
            <DateField label="updated" source="updated" showTime />
            <EditButton />
            <DeleteButton />
        </Datagrid>
    </List>
);

export const AlertEdit = (props) => (
    <Edit {...props}>
        <SimpleForm>
            <TextField source="crca" />
            <TextField source="crcs" />
            <TextField source="ip" />
            <TextField source="filesurvey" />
            <TextField source="role" />
            <TextField source="level" />
            <TextField source="lines"  style={{whiteSpace: 'pre-line'}}/>
            <RichTextInput source="comment" />
            <DateField label="created" source="created" showTime />
            <DateField label="updated" source="updated" showTime />
        </SimpleForm>
    </Edit>
);
