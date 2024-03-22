import React from 'react';
import { useMediaQuery } from '@mui/material';
import { List, Datagrid, TextField, DateField, EditButton,
    Edit, SimpleForm, Filter, TextInput, SelectInput, useRecordContext,
    RichTextField, SimpleList, Labeled
} from 'react-admin';
import { RichTextInput } from 'ra-input-rich-text';
import { roles } from './MyConfig';


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


const AlertFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Line" source="line" />
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} alwaysOn />
    </Filter>
);

const styles = {
    field: {
        display: 'inline-block', width: '250px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const AlertList = ({ classes, ...props }) => {
    const isSmall = useMediaQuery(theme => theme.breakpoints.down('sm'));
    return (
        <List filters={<AlertFilter />}  sort={{ field: 'updated', order: 'DESC' }} perPage={50} {...props}>
            {isSmall ? (
                <SimpleList
                    primaryText={record => `[${record.level}] ${record.search}`}
                    secondaryText={record => record.line}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            ) : (
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="level" />
                    <TextField source="line" />
                    <RichTextField source="comment" sx={ styles.field } stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
};

export const AlertEdit = (props) => (
    <Edit {...props}>
        <SimpleForm>
            <Labeled label="Search">
                <TextField source="search" />
            </Labeled>
            <Labeled label="Agent">
                <TextField source="ip" label="Agent" />
            </Labeled>
            <Labeled label="File survey">
                <TextField source="filesurvey" />
            </Labeled>
            <Labeled label="Role">
                <TextField source="role" />
            </Labeled>
            <Labeled label="Level">
                <ColoredTextField source="level" />
            </Labeled>
            <Labeled label="Line">
                <TextField source="line" sx={{width: '100%', whiteSpace: 'pre-line'}}/>
            </Labeled>
            <RichTextInput source="comment" />
            <Labeled label="Created">
                <DateField label="created" source="created" showTime />
            </Labeled>
            <Labeled label="Updated">
                <DateField label="updated" source="updated" showTime />
            </Labeled>
        </SimpleForm>
    </Edit>
);
