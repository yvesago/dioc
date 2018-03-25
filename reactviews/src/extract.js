import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField, 
    RichTextField, SelectInput, Filter, Responsive, SimpleList,
    BooleanInput, BooleanField
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

import { roles } from './MyConfig';

const actions = [
    { name: 'AddIP', id: 'AddIP' },
    { name: 'Delete', id: 'Delete' },
];

const colored = WrappedComponent => props => props.record.actions === 'Delete' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'warn' ? 
        <span style={{ color: 'orange' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);


const ExtractFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} allowEmpty alwaysOn />
    </Filter>
);


export const ExtractList = (props) => (
    <List filters={<ExtractFilter />} perPage={30} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => record.search}
                    secondaryText={record => `Role: ${record.role}, Active: ${record.active}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            }
            medium={
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="action" />
                    <BooleanField source="active" />
                    <RichTextField source="comment" stripTags elStyle={{width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}} />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            }
        />
    </List>
);


export const ExtractCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required} />
            <SelectInput source="role" choices={roles} allowEmpty />
            <SelectInput source="action" choices={actions} allowEmpty />
            <BooleanInput source="active" />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const ExtractEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required} />
            <SelectInput source="role" choices={roles} allowEmpty />
            <SelectInput source="action" choices={actions} allowEmpty />
            <BooleanInput source="active" />
            <RichTextInput source="comment" />
            <DateField label="Created" source="created" showTime />
            <DateField label="Updated" source="updated" showTime />
        </SimpleForm>
    </Edit>
);
