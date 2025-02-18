import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'category',
      label: 'Getting Started',
      items: [{ type: "autogenerated", dirName: "getting-started" }]
    },
    {
      type: 'category',
      label: 'User Guide',
      items: [{ type: "autogenerated", dirName: "user-guide" }]
    },
    {
      type: 'category',
      label: 'Hardware Acceleration',
      items: [{ type: "autogenerated", dirName: "hardware-acceleration" }]
    },
    {
      type: 'category',
      label: 'Application Customization',
      items: [{ type: "autogenerated", dirName: "app-customization" }]
    },
    {
      type: 'category',
      label: 'Advanced Topics',
      items: [{ type: "autogenerated", dirName: "advanced-topics" }]
    },
    {
      type: 'category',
      label: 'Developer Guide',
      items: [{ type: "autogenerated", dirName: "developer-guide" }]
    },
    {
      type: 'category',
      label: 'Deployment & Scaling',
      items: [{ type: "autogenerated", dirName:  "deployment-and-scaling" }]
    },
    {
      type: 'category',
      label: 'Integration & Extensibility',
      items: [{ type: "autogenerated", dirName: "integration-and-extensibility" }]
    },
    {
      type: 'category',
      label: 'Changelog & Releases',
      items: [{ type: "autogenerated", dirName: "changelog-and-releases" }]
    },
  ],
  apiSidebar: require("./docs/api/sidebar.js"),
};

export default sidebars;
