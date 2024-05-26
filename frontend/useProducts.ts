import { useEffect, useState } from 'react';
import axios from 'axios';

interface Product {
  id: string;
  name: string;
  price: number;
  // Any other product properties
}

interface UseProductOptions {
  initialSearch: string;
  initialSort: string;
  initialFilters: Record<string, string>;
}

const useProducts = ({ initialSearch, initialSort, initialFilters }: UseProductOptions) => {
  const [products, setProducts] = useState<Product[]>([]);
  const [search, setSearch] = useState(initialSearch);
  const [sort, setSort] = useState(initialSort);
  const [filters, setFilters] = useState(initialFilters);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<null | string>(null);

  // Fetch products
  const fetchProducts = async (search: string, sort: string, filters: Record<string, string>) => {
    setLoading(true);
    setError(null);
    
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('search', search);
      queryParams.append('sort', sort);
      Object.keys(filters).forEach(key => queryParams.append(key, filters[key]));

      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/products?${queryParams}`);
      setProducts(response.data);
    } catch (err: any) {
      setError(err.message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  // Update search
  const updateSearch = (search: string) => {
    setSearch(search);
  };

  // Update sort
  const updateSort = (sort: string) => {
    setSort(sort);
  };

  // Update filters
  const updateFilters = (filters: Record<string, string>) => {
    setFilters(filters);
  };

  // Effect for fetching products
  useEffect(() => {
    fetchProducts(search, sort, filters);
  }, [search, sort, filters]);

  return { products, loading, error, updateSearch, updateSort, updateFilters };
};

export default useProducts;