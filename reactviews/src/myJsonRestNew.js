import { fetchUtils } from 'react-admin';
import { stringify } from 'query-string';

//const apiUrl = 'https://my.api.com/';
//const httpClient = fetchUtils.fetchJson;

const filterQuery = value => {
    if (value && Object.keys(value).length) {
        return { _filters: JSON.stringify(value)};
    }
};

export default(apiUrl, httpClient = fetchUtils.fetchJson) => ({
    getList: async (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
            ...filterQuery(params.filter),
            _sortField: field,
            _sortDir: order,
            _start: (page - 1) * perPage,
            _end: page * perPage,
        };

        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const { json, headers } = await httpClient(url);
        if (resource === 'agents') {
            if ( json instanceof Array === true ) {
                json.map(x => { if (x.id === undefined ) { x.id = x.crca; } return x; });
            }
            else {
                if (json.id === undefined ) { json.id = json.crca; }
            }
        }

        return {
            data: json,
            total: parseInt(headers.get('x-total-count').split('/').pop(), 10),
        };
    },

    getOne: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const { json } = await httpClient(url);
        if (resource === 'agents') {
            if ( json instanceof Array === true ) {
                json.map(x => { if (x.id === undefined ) { x.id = x.crca; } return x; });
            }
            else {
                if (json.id === undefined ) { json.id = json.crca; }
            }
        }
        return { data: json };
    },

    getMany: async (resource, params) => {
        const query = {
            filter: JSON.stringify({ ids: params.ids }),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const { json } = await httpClient(url);
        if (resource === 'agents') {
            if ( json instanceof Array === true ) {
                json.map(x => { if (x.id === undefined ) { x.id = x.crca; } return x; });
            }
            else {
                if (json.id === undefined ) { json.id = json.crca; }
            }
        }
        return { data: json };
    },

    getManyReference: async (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        /*const query = {
            sort: JSON.stringify([field, order]),
            range: JSON.stringify([(page - 1) * perPage, page * perPage - 1]),
            filter: JSON.stringify({
                ...params.filter,
                [params.target]: params.id,
            }),
        };*/
        const query = {
            ...filterQuery(params.filter),
            [params.target]: params.id,
            _sortField: field,
            _sortDir: order,
            _start: (page - 1) * perPage,
            _end: page * perPage,
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const { json, headers } = await httpClient(url);
        if (resource === 'agents') {
            if ( json instanceof Array === true ) {
                json.map(x => { if (x.id === undefined ) { x.id = x.crca; } return x; });
            }
            else {
                if (json.id === undefined ) { json.id = json.crca; }
            }
        }
        return {
            data: json,
            total: parseInt(headers.get('x-total-count').split('/').pop(), 10),
        };
    },

    create: async (resource, params) => {
        const { json } = await httpClient(`${apiUrl}/${resource}`, {
            method: 'POST',
            body: JSON.stringify(params.data),
        });
        return { data: json };
    },

    update: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const { json } = await httpClient(url, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        });
        if (resource === 'agents') {
            if ( json instanceof Array === true ) {
                json.map(x => { if (x.id === undefined ) { x.id = x.crca; } return x; });
            }
            else {
                if (json.id === undefined ) { json.id = json.crca; }
            }
        }
        return { data: json };
    },

    updateMany: async (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids}),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const { json } = await httpClient(url, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        });
        if (resource === 'agents') {
            if ( json instanceof Array === true ) {
                json.map(x => { if (x.id === undefined ) { x.id = x.crca; } return x; });
            }
            else {
                if (json.id === undefined ) { json.id = json.crca; }
            }
        }
        return { data: json };
    },

    delete: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const { json } = await httpClient(url, {
            method: 'DELETE',
        });
        return { data: json };
    },

    deleteMany: async (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids}),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const { json } = await httpClient(url, {
            method: 'DELETE',
            body: JSON.stringify(params.data),
        });
        return { data: json };
    },
});

