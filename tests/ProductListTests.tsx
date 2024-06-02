import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; 
import ProductList from './ProductList';
import axiosMock from 'axios';

process.env.REACT_APP_API_URL = 'http://your-api-url.com';

const mockProductData = [
  { id: 1, name: 'Product 1', price: 10, category: 'Category 1' },
  { id: 2, name: 'Product 2', price: 20, category: 'Category 2' },
  { id: 3, name: 'Product 3', price: 30, category: 'Category 1' },
];

axiosMock.get.mockResolvedValue({ data: mockProductData });

describe('ProductList Component', () => {
  it('renders product data fetched from API', async () => {
    const { getByText } = render(<ProductList />);
    
    await waitFor(() => {
      expect(getByText('Product 1')).toBeInTheDocument();
      expect(getByText('Product 2')).toBeInTheDocument();
      expect(getByText('Product 3')).toBeInTheDocument();
    });
  });

  it('filters products by category correctly', async () => {
    const { getByText, getByTestId } = render(<ProductList />);
    const filterInput = getByTestId('filter-input');

    fireEvent.change(filterInput, { target: { value: 'Category 1' } });

    await waitFor(() => {
      expect(getByText('Product 1')).toBeInTheDocument();
      expect(getByText('Product 3')).toBeInTheDocument();
      expect(() => getByText('Product 2')).toThrow();
    });
  });

  it('sorts products by price correctly', async () => {
    const { getByText, getByTestId } = render(<ProductList />);
    const sortButton = getByTestId('sort-button-price');

    fireEvent.click(sortButton);

    const sortedPrices = mockProductData
      .sort((a, b) => a.price - b.price)
      .map(product => product.price);
    
    await waitFor(() => {
      sortedPrices.forEach((price, index) => {
        expect(getByText(`Product ${index + 1}`)).toBeInTheDocument();
      });
    });
  });
});