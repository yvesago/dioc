import React from 'react';
import { Edit, SimpleForm  } from 'admin-on-rest';
import RichTextInput from 'aor-rich-text-input';

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
