import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; 
import ProductList from './ProductList';
import axiosMock from 'axios';

process.env.REACT_APP_API_URL = 'http://your-api-url.com';

const mockProductData = [
  { id: 1, name: 'Product 1', price: 10, category: 'Category 1', dynamicPrice: () => 9 },
  { id: 2, name: 'Product 2', price: 20, category: 'Category 2', dynamicPrice: () => 18 },
  { id: 3, name: 'Product 3', price: 30, category: 'Category 1', dynamicPrice: () => 27 },
];

axiosMock.get.mockResolvedValue({ data: mockProductData });

describe('ProductList Component', () => {
  it('updates product prices dynamically based on user interaction', async () => {