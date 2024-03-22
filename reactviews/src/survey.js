import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField, useRecordContext,
    RichTextField, SelectInput, Filter, SimpleList, Labeled
} from 'react-admin';
import { useMediaQuery } from '@mui/material';
import { RichTextInput } from 'ra-input-rich-text';
import { roles } from './MyConfig';

const levels = [
    { name: 'Critic', id: 'critic' },
    { name: 'Warn', id: 'warn' },
];


const ColoredTextField = (props) => {
    const record = useRecordContext();
    return (
        <TextField
            sx={{ color: record.level === 'critic' ? 'red' : 'orange' }}
            {...props}
        />
    );
};
ColoredTextField.defaultProps = TextField.defaultProps;


const SurveyFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} alwaysOn />
    </Filter>
);

const styles = {
    field: {
        display: 'inline-block', width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const SurveyList = ({ classes, ...props }) => {
    const isSmall = useMediaQuery(
        theme => theme.breakpoints.down('sm'),
        { noSsr: true }
    );
    return (
        <List bulkActionButtons={false} filters={<SurveyFilter />} sort={{ field: 'updated', order: 'DESC' }} perPage={50} {...props}>
            {isSmall ? (
                <SimpleList
                    primaryText={record => record.search}
                    secondaryText={record => `Role: ${record.role}, Level: ${record.level}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            ) : (
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="level" />
                    <RichTextField source="comment" sx={ styles.field } stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
};


export const SurveyCreate = (props) => (
    <Create {...props}>
        <SimpleForm redirect="list">
            <TextInput source="search" validate={required()} />
            <SelectInput source="level" choices={levels} />
            <SelectInput source="role" choices={roles} />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const SurveyEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required()} />
            <SelectInput source="role" choices={roles} />
            <SelectInput source="level" choices={levels} />
            <RichTextInput source="comment" />
            <Labeled label="Created">
                <DateField source="created" />
            </Labeled>
            <Labeled label="Updated">
                <DateField source="updated" />
            </Labeled>
            <TextInput source="crcs" validate={required()} disabled />
        </SimpleForm>
    </Edit>
);
