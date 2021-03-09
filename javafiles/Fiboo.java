public class Fiboo{

    private static long dfs(long n){
    	if(n <= 1)	{return n;}
    	return dfs(n-1) + dfs(n-2);
    }

    public static void main(String[] args){
        long m = dfs(30);
        System.out.println(m);
    }
}