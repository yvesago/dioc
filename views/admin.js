
var myApp = angular.module('myApp', ['restangular','ng-admin']);
myApp.config(['NgAdminConfigurationProvider','RestangularProvider', function (nga,rp) {

    rp.setDefaultHeaders({"X-MyToken":'basic'}); // XXX fix token

    // create an admin application
    var admin = nga.application('My REST Admin')
      .baseApiUrl('http://localhost:8080/admin/api/v1/'); // XXX fix main API endpoint

    // create a survey entity
    var survey = nga.entity('surveys');
    survey.listView().fields([
        nga.field('search').isDetailLink(true),
        nga.field('role'),
        nga.field('level').cssClasses(['red']).template('<span ng-class="{ \'orange\': value === \'warn\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
        nga.field('crcs'),
        nga.field('comment','wysiwyg'),
        nga.field('created','datetime'),
        nga.field('updated','datetime'),
    ]).filters([
            nga.field('search')
            .label('Search')
            .pinned(true),
            nga.field('comment')
            .label('Comment')
            ]);
    survey.creationView().fields([
        nga.field('search'),
        nga.field('level', 'choice').choices([
                       { label: 'Warn', value: 'warn' },
                       { label: 'Critic', value: 'critic' },
                       { label: 'None', value: '' }
                   ]),
        nga.field('role', 'choice').choices([ // XXX fix your categorys
                       { label: 'Squid', value: 'squid' },
                       { label: 'Radius', value: 'radius' },
                       { label: 'Web', value: 'web' },
                       { label: 'DNS', value: 'dns' },
                       { label: 'Honeypot', value: 'honeypot' },
                       { label: 'NetFlow', value: 'netflow' },
                       { label: 'Auth', value: 'auth' },
                       { label: 'SMTP', value: 'smtp' },
                       { label: 'Mail', value: 'mail' },
                       { label: 'Test', value: 'test' },
                       { label: 'None', value: '' }
                   ]),
        nga.field('comment','wysiwyg'),
    ]);
    // use the same fields for the editionView as for the creationView
    survey.editionView().fields(survey.creationView().fields());
    admin.addEntity(survey);

    var agent = nga.entity('agents').identifier(nga.field('crca'));
    agent.listView().fields([
       nga.field('crca').isDetailLink(true),
       nga.field('ip'),
       nga.field('status').cssClasses(['red','green']).template('<span ng-class="{ \'red\': value === \'OffLine\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
       nga.field('role'),
       nga.field('filesurvey'),
       nga.field('lines','text'),
       nga.field('comment','text'),
       nga.field('cmd'),
       nga.field('created','datetime'),
       nga.field('updated','datetime'),
    ]).filters([
            nga.field('ip')
            .label('IP')
            .pinned(true),
            ]);
    agent.creationView().fields([
       nga.field('crca'),
       nga.field('ip'),
       nga.field('filesurvey'),
       nga.field('lines','text'),
       nga.field('comment','text'),
       nga.field('role', 'choice').choices([ // XXX fix your categorys
                       { label: 'Squid', value: 'squid' },
                       { label: 'Radius', value: 'radius' },
                       { label: 'Web', value: 'web' },
                       { label: 'DNS', value: 'dns' },
                       { label: 'Honeypot', value: 'honeypot' },
                       { label: 'NetFlow', value: 'netflow' },
                       { label: 'Auth', value: 'auth' },
                       { label: 'SMTP', value: 'smtp' },
                       { label: 'Mail', value: 'mail' },
                       { label: 'Test', value: 'test' },
                       { label: 'None', value: '' }
                   ]),
       nga.field('cmd', 'choice').choices([
                       { label: 'Send Lines', value: 'SendLines' },
                       { label: 'Search Full File', value: 'FullSearch' },
                       { label: 'STOP', value: 'STOP' },
                       { label: 'None', value: '' }
                   ]),
    ]);
    agent.editionView().fields(agent.creationView().fields());
    admin.addEntity(agent);

    var alerte = nga.entity('alertes');
    alerte.listView().fields([
        nga.field('search').isDetailLink(true),
        nga.field('role'),
        nga.field('ip'),
        nga.field('filesurvey'),
        nga.field('level').cssClasses(['red']).template('<span ng-class="{ \'orange\': value === \'warn\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
        nga.field('line'),
        nga.field('comment','text'),
        nga.field('created','datetime'),
        nga.field('updated','datetime'),
    ]).filters([
            nga.field('line')
            .label('Line')
            .pinned(true)]);
    alerte.creationView().fields([
        nga.field('crca'),
        nga.field('crcs'),
        nga.field('line'),
        nga.field('comment','text'),
        //nga.field('comment','wysiwyg'),
    ]);
    alerte.editionView().fields(alerte.creationView().fields());
    admin.addEntity(alerte);
    // attach the admin application to the DOM and execute it
    nga.configure(admin);
}]);
