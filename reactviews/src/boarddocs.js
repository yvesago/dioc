import React from 'react';
import { Edit, SimpleForm  } from 'react-admin';
import RichTextInput from 'ra-input-rich-text';

export const BoardEdit = (props) => (
    <Edit title="Docs" {...props}>
        <SimpleForm redirect="show">
            <RichTextInput
                source="docs"
                label="Main Doc"
            />
            <RichTextInput
                source="docagents"
                label="Agents"
            />
            <RichTextInput
                source="docsearchs"
                label="Searchs"
            />
        </SimpleForm>
    </Edit>
);
