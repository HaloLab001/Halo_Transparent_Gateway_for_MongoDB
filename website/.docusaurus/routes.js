import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/markdown-page/',
    component: ComponentCreator('/markdown-page/', '3c9'),
    exact: true
  },
  {
    path: '/',
    component: ComponentCreator('/', '2bc'),
    exact: true
  },
  {
    path: '/',
    component: ComponentCreator('/', '80d'),
    routes: [
      {
        path: '/Basic_operations/',
        component: ComponentCreator('/Basic_operations/', '4d6'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/Basic_operations/create/',
        component: ComponentCreator('/Basic_operations/create/', '528'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/Basic_operations/delete/',
        component: ComponentCreator('/Basic_operations/delete/', '0bb'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/Basic_operations/read/',
        component: ComponentCreator('/Basic_operations/read/', '89d'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/Basic_operations/update/',
        component: ComponentCreator('/Basic_operations/update/', 'e23'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/category/basic-crud-operations/',
        component: ComponentCreator('/category/basic-crud-operations/', 'a63'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/category/quickstart/',
        component: ComponentCreator('/category/quickstart/', '9d1'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/contributing/',
        component: ComponentCreator('/contributing/', 'b7d'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/diff/',
        component: ComponentCreator('/diff/', 'b22'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/intro/',
        component: ComponentCreator('/intro/', '785'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/quickstart_guide/debian/',
        component: ComponentCreator('/quickstart_guide/debian/', '000'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/quickstart_guide/docker/',
        component: ComponentCreator('/quickstart_guide/docker/', '585'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/quickstart_guide/macos/',
        component: ComponentCreator('/quickstart_guide/macos/', 'c31'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/quickstart_guide/rpm/',
        component: ComponentCreator('/quickstart_guide/rpm/', '5e5'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/quickstart_guide/windows/',
        component: ComponentCreator('/quickstart_guide/windows/', '652'),
        exact: true,
        sidebar: "sidebar"
      },
      {
        path: '/understanding_ferretdb/',
        component: ComponentCreator('/understanding_ferretdb/', '076'),
        exact: true,
        sidebar: "sidebar"
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
