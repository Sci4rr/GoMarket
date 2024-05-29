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

  const fetchProducts = async (search: string, sort: string, filters: Record<string, string>) => {
    setLoading(true);
    setError(null);

    try {
      const queryParams = new URLSearchParams();
      queryParams.append('search', search);
      queryParams.append('sort', sort);
      Object.keys(filters).forEach(key => queryParams.append(key, filters[key]));

      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/products?${queryParams}`);
      if(response.data && Array.isArray(response.data)) {
        setProducts(response.data);
      } else {
        throw new Error("Unexpected response format");
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        if (err.response) {
          setError(`Server responded with status code ${err.response.status}: ${err.message}`);
        } else if (err.request) {
          setError("The request was made but no response was received");
        } else {
          setError("An error occurred while setting up the request");
        }
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setLoading(false);
    }
  };

  const updateSearch = (newSearch: string) => {
    setSearch(newSearch);
  };

  const updateSort = (newSort: string) => {
    setSort(newSort);
  };

  const updateFilters = (newFilters: Record<string, string>) => {
    setFilters(newFilters);
  };

  useEffect(() => {
    fetchProducts(search, sort, filters);
  }, [search, sort, filters]);

  return { products, loading, error, updateSearch, updateSort, updateFilters };
};

export default useProducts;