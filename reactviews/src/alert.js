import React from 'react';
import { List, Datagrid, TextField, DateField, EditButton,
    Edit, SimpleForm
} from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

export const AlertList = (props) => (
    <List {...props}>
        <Datagrid>
            <TextField source="search" />
            <TextField source="role" />
            <TextField source="ip" />
            <TextField source="filesurvey" />
            <TextField source="level" />
            <TextField source="line" />
            <TextField source="comment" />
            <DateField label="updated" source="updated" showTime />
            <EditButton />
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
